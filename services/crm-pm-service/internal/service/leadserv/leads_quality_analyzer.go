package leadserv

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
)

type QualityIndicator string

const (
	LengthIndicator      QualityIndicator = "length"
	StructureIndicator   QualityIndicator = "structure"
	SpecificityIndicator QualityIndicator = "specificity"
	UrgencyIndicator     QualityIndicator = "urgency"
	NumbersIndicator     QualityIndicator = "numbers"
)

type LocalLeadAnalyzer struct {
	//** Peso de cada indicador (configurable)
	weights  map[QualityIndicator]float64
	keywords *keywords
}

type keywords struct {
	ActionVerbs     []string `json:"action_verbs"`
	TechnicalWords  []string `json:"technical_words"`
	UrgencyKeywords []string `json:"urgency_keywords"`
	TimeKeywords    []string `json:"time_keywords"`
}

func NewLocalLeadAnalyzer() *LocalLeadAnalyzer {
	keywordsPath := os.Getenv("KEYWORDS_PATH")
	if keywordsPath == "" {
		keywordsPath = "./data/keywords.json"
	}

	kw, err := loadKeywords(keywordsPath)
	if err != nil {
		panic(fmt.Sprintf("FATAL: failed_to_load_keywords_from %s: %v", keywordsPath, err))
	}

	return &LocalLeadAnalyzer{
		weights: map[QualityIndicator]float64{
			LengthIndicator:      0.15, // 15% del score
			StructureIndicator:   0.25, //25%
			SpecificityIndicator: 0.35, // 35%
			UrgencyIndicator:     0.15, // 15%
			NumbersIndicator:     0.10, // 10%
		},
		keywords: kw,
	}
}

func loadKeywords(filePath string) (*keywords, error) {
	file, err := os.Open(filePath) //nolint:gosec // G304: path from env config (KEYWORDS_PATH), not user input
	if err != nil {
		return nil, fmt.Errorf("failed_to_open_keywords_file: %w", err)
	}
	defer file.Close() //nolint:errcheck // read-only file; Close error is not actionable

	var kw keywords
	if err := json.NewDecoder(file).Decode(&kw); err != nil {
		return nil, fmt.Errorf("failed_to_parse_keywords_file: %w", err)
	}

	return &kw, nil
}

func (a *LocalLeadAnalyzer) Analyze(message string) *domain.QualityAnalisis {
	analysis := &domain.QualityAnalisis{
		Indicators:    make(map[string]float64),
		FoundConcepts: []string{},
	}

	analysis.Indicators[string(LengthIndicator)] = a.analyzeLength(message)
	analysis.Indicators[string(StructureIndicator)] = a.analyzeStructure(message)
	analysis.Indicators[string(SpecificityIndicator)] = a.analyzeSpecificity(message)
	analysis.Indicators[string(UrgencyIndicator)] = a.analyzeUrgency(message)
	analysis.Indicators[string(NumbersIndicator)] = a.analyzeNumbers(message)

	weightedScore := 0.0
	for indicator, score := range analysis.Indicators {
		weight := a.weights[QualityIndicator(indicator)]
		weightedScore += score * weight
	}

	analysis.Score = weightedScore * 100

	switch {
	case analysis.Score >= 70:
		analysis.QA = domain.High
	case analysis.Score >= 40:
		analysis.QA = domain.Medium
	default:
		analysis.QA = domain.Low
	}

	analysis.FoundConcepts = a.extractConcepts(message)

	return analysis
}

// ** ========== INDICADORES ==========

// mensajes muy cortos = contexto insuficiente
func (a *LocalLeadAnalyzer) analyzeLength(message string) float64 {
	wordCount := len(strings.Fields(message))

	switch {
	case wordCount < 20:
		return 0.0 // Demasiado corto
	case wordCount < 50:
		return 0.4 // Corto
	case wordCount < 100:
		return 0.7
	case wordCount < 300:
		return 1.0
	default:
		return 1.0
	}
}

// preguntas, puntuacion, parrafos = claridad
func (a *LocalLeadAnalyzer) analyzeStructure(message string) float64 {
	score := 0.0

	// Sub-indicador: preguntas (indica pensamiento estructurado)
	questionCount := strings.Count(message, "?")
	if questionCount > 0 {
		score += 0.3
	}

	// Sub-indicador: multiples oraciones
	sentenceCount := strings.Count(message, ".") + questionCount
	if sentenceCount > 2 {
		score += 0.3
	}

	// Sub-indicador: parrafos (saltos de linea)
	//! Verificar si funciona con \n
	paragraphs := strings.Count(message, "\n") + 1
	if paragraphs > 1 {
		score += 0.2
	}

	if !a.isAllCaps(message) && a.hasProperCapitalization(message) {
		score += 0.2
	}

	return math.Min(score, 1.0)
}

func (a *LocalLeadAnalyzer) analyzeSpecificity(message string) float64 {
	score := 0.0

	lower := strings.ToLower(message)

	// Sub-indicador: numeros (presupuesto, timeline, cantidades)
	if a.hasNumbers(message) {
		score += 0.3
	}

	// Sub-indicador: verbos de accion
	if a.containsAny(lower, a.keywords.ActionVerbs) {
		score += 0.25
	}

	// Sub-indicador: palabras tecnicas
	techCount := a.countMatches(lower, a.keywords.TechnicalWords)
	if techCount > 0 {
		score += 0.2 * math.Min(float64(techCount)/3.0, 1.0)
	}

	// Sub-indicador: plazos/timeline
	if a.containsAny(lower, a.keywords.TimeKeywords) {
		score += 0.25
	}

	return math.Min(score, 1.0)
}

// indicador de cliente serio vs curiosidad
func (a *LocalLeadAnalyzer) analyzeUrgency(message string) float64 {
	lower := strings.ToLower(message)

	if a.containsAny(lower, a.keywords.UrgencyKeywords) {
		return 0.8
	}

	return 0.2
}

// presupuesto, cantidades, metricas
func (a *LocalLeadAnalyzer) analyzeNumbers(message string) float64 {
	if !a.hasNumbers(message) {
		return 0.0
	}

	moneyPatterns := regexp.MustCompile(
		`\$\d+|€\d+|\d+\s*(k|m|usd|eur|pesos|dolares|euros)`,
	)
	if moneyPatterns.MatchString(strings.ToLower(message)) {
		return 1.0
	}

	return 0.6
}

// ** ========== HELPERS ==========

func (a *LocalLeadAnalyzer) isAllCaps(text string) bool {
	lower := strings.ToLower(text)
	upper := strings.ToUpper(text)
	return len(text) > 0 && text == upper && lower != upper
}

func (a *LocalLeadAnalyzer) hasProperCapitalization(text string) bool {
	words := strings.Fields(text)
	if len(words) == 0 {
		return false
	}

	properCount := 0
	for _, word := range words {
		if len(word) > 0 && unicode.IsUpper(rune(word[0])) {
			properCount++
		}
	}

	return float64(properCount)/float64(len(words)) > 0.3
}

func (a *LocalLeadAnalyzer) hasNumbers(text string) bool {
	for _, r := range text {
		if unicode.IsDigit(r) {
			return true
		}
	}

	return false
}

func (a *LocalLeadAnalyzer) containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}

	return false
}

func (a *LocalLeadAnalyzer) countMatches(text string, keywords []string) int {
	count := 0

	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			count++
		}
	}
	return count
}

// detecta conceptos clave para logging
func (a *LocalLeadAnalyzer) extractConcepts(message string) []string {
	lower := strings.ToLower(message)
	concepts := []string{}
	seen := make(map[string]bool)

	allKeywords := append(
		append(
			append(a.keywords.ActionVerbs, a.keywords.TechnicalWords...),
			a.keywords.UrgencyKeywords...,
		),
		a.keywords.TimeKeywords...,
	)

	for _, keyword := range allKeywords {
		if strings.Contains(lower, keyword) && !seen[keyword] {
			concepts = append(concepts, keyword)
			seen[keyword] = true
		}
	}

	return concepts
}
