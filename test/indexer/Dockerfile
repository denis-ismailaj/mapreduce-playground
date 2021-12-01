FROM mr-playground-test as base

COPY ./test/indexer/. .


FROM base as build

COPY . .

RUN go build "$RACE" -buildmode=plugin indexer.go


FROM base

COPY --from=build /app/indexer.so .

ENV DATA_DIR="testdata"

RUN chmod +x indexer.sh

ENTRYPOINT ["bash", "-c", "./indexer.sh"]
