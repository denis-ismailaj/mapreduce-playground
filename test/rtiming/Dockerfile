FROM mr-playground-test as base

COPY ./test/rtiming/. .


FROM base as build

COPY . .

RUN go build "$RACE" -buildmode=plugin rtiming.go


FROM base

COPY --from=build /app/rtiming.so .

ENV DATA_DIR="testdata"

RUN chmod +x rtiming.sh

ENTRYPOINT ["bash", "-c", "./rtiming.sh"]
