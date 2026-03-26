FROM heroiclabs/nakama-pluginbuilder:3.38.0 AS builder

WORKDIR /backend
COPY . .

RUN go build -C=./plugin -o ../plugin.so -buildmode=plugin -trimpath

FROM heroiclabs/nakama:3.38.0

COPY --from=builder /backend/plugin.so /nakama/data/modules
COPY --from=builder /backend/local.yml /nakama/data/
COPY --from=builder /backend/*.json /nakama/data/modules
