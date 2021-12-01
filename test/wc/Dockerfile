FROM mr-playground-test as base

COPY ./test/wc/. .


FROM base as build

COPY . .

RUN go build "$RACE" -buildmode=plugin wc.go


FROM base

COPY --from=build /app/wc.so .

ENV DATA_DIR="testdata"

RUN chmod +x wc.sh

ENTRYPOINT ["bash", "-c", "./wc.sh"]
