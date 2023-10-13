FROM debian:stable-slim

COPY personal-finance-api /bin/personal-finance-api

CMD ["/bin/personal-finance-api"]
