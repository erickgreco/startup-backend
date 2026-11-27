package fonts

import _ "embed"

// ! Por ahora se usa asi para tests

//go:embed Inter/static/Inter_18pt-Light.ttf
var InterLight []byte

//go:embed Inter/static/Inter_18pt-Regular.ttf
var InterRegular []byte

//go:embed Inter/static/Inter_18pt-SemiBold.ttf
var InterSemiBold []byte
