
# creditcard

Currect repo is implementation of the credit card task on the alem.school platform. It fulfills all the conditions that are inside it

## Usage

Download repo, and then run code using command below

```bash
go run ./cmd/app
```

Programm implements next functions

`Generate` - The generate feature creates possible credit card numbers by replacing asterisks (*) with digits.


```bash
 go run ./cmd/app generate --pick "440043018030****"
```

`validate` - The validate feature checks if a credit card number is valid using Luhn's Algorithm.

```bash
go run ./cmd/app validate "4400430180300003"
OK
```

`information` - The information feature provides details about the card based on data in brands.txt and issuers.txt.

```bash
go run ./cmd/app information --brands=brands.txt --issuers=issuers.txt "4400430180300003"
4400430180300003
Correct: yes
Card Brand: VISA
Card Issuer: Kaspi Gold
```

`issue` - The issue feature generates a random valid credit card number for a specified brand and issuer.

```bash
go run ./cmd/app issue --brands=brands.txt --issuers=issuers.txt --brand=VISA --issuer="Kaspi Gold"
```
