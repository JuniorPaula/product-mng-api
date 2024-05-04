# Simples API de gerenciamento de produtos

Essa é uma API simples desenvolvida em `Golang` utilizando as boas práticas e mantendo o código modularizado.
O Sistema conta com autenticação `JWT`, uma simples persistência utilizando `SQLite3` e testes automatizados.

## Referências
[Chi](github.com/go-chi/chi/v5)
[Mock](github.com/golang/mock)
[Google uuid](github.com/google/uuid)
[Viper](github.com/spf13/viper)
[Testfy](github.com/stretchr/testify)
[SQLite3](gorm.io/driver/sqlite)
[GORM](gorm.io/gorm)


### Funcionalidades

- Cadastrar produto
- Atualizar produto
- Listar produto
- Deletar produto
- Sistema de login com autenticação utilizando jwt.
- Cadastrar usuário (apenas admin)
- Atualizar usuário (apenas admin)
- Deletar usuário (apenas admin)

### Executar
```bash
go build -o main ./cmd/server/main.go
./main
```