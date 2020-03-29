FROM golang
COPY ./bot-manager .
RUN chmod +x ./bot-manager
RUN ./bot-manager
