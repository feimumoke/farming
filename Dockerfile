FROM debian

#eg docker build -t xxx --build-arg PROJECT_NAME=server-farmer .
ARG PROJECT_NAME
ARG GRPC_PORT
ARG HTTP_PORT
ENV PROJECT_NAME=${PROJECT_NAME:-demo_project} \
    GRPC_PORT=${GRPC_PORT:-28086} \
    HTTP_PORT=${HTTP_PORT:-28088}

COPY target /usr/local/services/

EXPOSE $GRPC_PORT
EXPOSE $HTTP_PORT

WORKDIR /usr/local/services/$PROJECT_NAME
RUN chmod +x ./bin/start.sh && sync

CMD ["bin/start.sh", "$PROJECT_NAME"]