FROM mr-playground-test as base

COPY ./test/jobcount/. .


FROM base as build

COPY . .

RUN go build "$RACE" -buildmode=plugin jobcount.go


FROM base

COPY --from=build /app/jobcount.so .

ENV DATA_DIR="testdata"

RUN chmod +x jobcount.sh

ENTRYPOINT ["bash", "-c", "./jobcount.sh"]
