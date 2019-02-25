FROM golang:1.11.5

LABEL version="1.0.0"
LABEL repository="https://github.com/actions/go"
LABEL homepage="https://github.com/actions"
LABEL maintainer="GitHub Actions <support+actions@github.com>"

LABEL com.github.actions.name="GitHub Action for Go"
LABEL com.github.actions.description="Go wrapper"
LABEL com.github.actions.icon="package"
LABEL com.github.actions.color="blue"

ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]