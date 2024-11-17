# Etapa de build
FROM golang:1.23.2 as builder

# Definir diretório de trabalho
WORKDIR /app

# Copiar o código da aplicação
COPY . .

# Instalar pacotes necessários, incluindo certificados
RUN apt-get update && apt-get install -y ca-certificates

# Garantir que os certificados sejam atualizados
RUN update-ca-certificates

# Construir a aplicação com o Go (desabilitando CGO para garantir binário estático)
RUN GOOS=linux CGO_ENABLED=0 go build -o server ./cmd/server

# Etapa de produção
FROM alpine:latest
# Usar uma imagem leve, mas com suporte a certificados

# Instalar certificados na imagem de produção (caso o servidor precise fazer requisições HTTPS)
RUN apk add --no-cache ca-certificates

# Copiar o binário da etapa de build
COPY --from=builder /app/server /server

# Definir o comando padrão para rodar o binário
CMD ["/server"]
