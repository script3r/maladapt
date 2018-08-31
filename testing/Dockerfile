FROM mongo

COPY ./start.sh /etc/mongo/

COPY ./mongodb.pem /etc/ssl/
COPY ./mongo.yaml /etc/mongo/
COPY ./configure.js /etc/mongo/

EXPOSE 27017

CMD ["./etc/mongo/start.sh"]
