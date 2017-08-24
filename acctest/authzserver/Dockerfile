FROM alpine:3.6

ENTRYPOINT ["/docker-entrypoint.sh"]
ADD server.crt /
ADD server.key /
ADD docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh
ADD authful.yml /
ADD authful-linux-386 /authful