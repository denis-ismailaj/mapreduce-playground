FROM mr-playground-test as base

COPY ./test/crash/. .


FROM base as build

COPY . .

RUN go build "$RACE" -buildmode=plugin crash.go

RUN go build "$RACE" -buildmode=plugin nocrash.go


FROM base

RUN apt-get update && apt-get install -y netcat

COPY --from=build /app/crash.so .

COPY --from=build /app/nocrash.so .

ENV DATA_DIR="testdata"

RUN chmod +x crash.sh

ENTRYPOINT ["bash", "-c", "./crash.sh"]
