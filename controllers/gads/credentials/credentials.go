package gads

import (
	"os"
)

var (
	R_T    = os.Getenv("R_T")    // refresh_token recuperado do seu oauth2 flow
	C_I    = os.Getenv("C_I")    // client ID do oauth gerado no google cloud
	C_S    = os.Getenv("C_S")    // client secreet do oauth gerado no google cloud
	D_T    = os.Getenv("D_T")    // developer token gerado na sua conta MCC e aprovado pelo google ads para acesso básico
	MCC_ID = os.Getenv("MCC_ID") // ID da sua MCC sem os traços "000-000-0000"
)
