
var ApiKey string
var sendMessageUrl string
var apiURL string
var OriginPhoneNumber string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}
	ApiKey = os.Getenv("WHATSGW_API_KEY")
	sendMessageUrl = os.Getenv("WHATSGW_SEND_MESSAGE")
	apiURL = os.Getenv("CONFIRM_CONSULTA")
	OriginPhoneNumber = os.Getenv("WHATSGW_PHONE_NUMBER")
}
