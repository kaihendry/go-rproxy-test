FROM scratch
ADD bbc /bbc
ENTRYPOINT ["/bbc"]
EXPOSE 9000
