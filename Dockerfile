FROM alpine
ADD router /router
ADD filebeat.yml /filebeat.yml
ENTRYPOINT [ "/router" ]
