# GoTraining

Run App: go run .cmd\CliApp -message="(message string)" -userID="(userID int)"
Exit TaskApp: CTRL + C

curl -X POST http://localhost:8080/messages   -H "Content-Type: application/json"   -d '{"message":"(message string)","userID":(userID int)}'