FROM mr-playground-test as base

COPY ./test/mtiming/. .


FROM base as build

COPY . .

RUN go build "$RACE" -buildmode=plugin mtiming.go


FROM base

COPY --from=build /app/mtiming.so .

ENV DATA_DIR="testdata"

RUN chmod +x mtiming.sh

ENTRYPOINT ["bash", "-c", "./mtiming.sh"]
