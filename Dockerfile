FROM scratch

COPY reindexer_exporter /

EXPOSE      9451

CMD ["/reindexer_exporter"]
