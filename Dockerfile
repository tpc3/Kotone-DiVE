FROM scratch
COPY ./bin/linux/amd64/Kotone-DiVE config.yaml /Kotone-DiVE/
ENTRYPOINT [ "Kotone-DiVE/Kotone-DiVE" ]