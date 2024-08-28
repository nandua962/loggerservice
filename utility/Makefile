server:
	go run main.go

test:
	go test -cover -v ./...

mock:
	mockgen -package mock -destination internal/repo/mock/language.go utility/internal/repo LanguageRepoImply
	mockgen -package=mock -destination=./internal/repo/mock/currency.go  utility/internal/repo CurrencyRepoImply
	mockgen -package=mock -destination=./internal/repo/mock/country.go  utility/internal/repo CountryRepoImply
