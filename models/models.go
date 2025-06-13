package models

type AppConfig struct {
	Redis                    RedisConfig  
	TPSConfig                map[string]int 
	Ports                    []int          
	MaxRequestsPerConnection int            
	TotalConnections         int            
	Server                   struct {
		Host string
		Port string
	}
	RedisExpiration             int
}

type TenantData struct{
	ID int
	Name string
	Phone string
	Status 	string
	ServerAck string 
}

type RedisConfig struct {
	Host        string
	RedisName           string
	RedisPassword       string
	DB                  int
}

type GenerateRequest struct {
	Count int `json:"count"`
}