FROM mr-playground-test as base

COPY ./test/early_exit/. .


FROM base as build

COPY . .

RUN go build "$RACE" -buildmode=plugin early_exit.go


FROM base

COPY --from=build /app/early_exit.so .

ENV DATA_DIR="testdata"

RUN chmod +x early_exit.sh

ENTRYPOINT ["bash", "-c", "./early_exit.sh"]
