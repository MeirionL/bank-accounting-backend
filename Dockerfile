FROM debian:stable-slim

ENV PORT=8080
ENV JWT_SECRET=qRcrjPQ989usFYG/O2fcBXASlGNisXI4v0+9bWNwgXCNXjYPKww4d6j93faBB8poxi9K7QzkWLSC8UGtvgNzOw==
ENV DB_URL=postgres://postgres:ElephantsFTW@localhost:5432/pfa-data?sslmode=disable

COPY out /bin/out

CMD ["/bin/out"]
