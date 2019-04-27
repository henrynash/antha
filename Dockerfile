FROM eu.gcr.io/antha-images/golang:1.12.4-build

ARG COMMIT_SHA
ADD .netrc /
RUN mv /.netrc $HOME/.netrc || true
RUN mkdir /antha
WORKDIR /antha
RUN set -ex && go mod init antha && go get github.com/antha-lang/antha@$COMMIT_SHA && go mod download
RUN set -ex && go install github.com/antha-lang/antha/cmd/...
RUN set -ex && go test -c github.com/antha-lang/antha/cmd/elements
RUN set -ex && go test github.com/antha-lang/antha/...
COPY scripts/. /antha/.
RUN rm $HOME/.netrc

# These are for the gitlab CI for elements:
ONBUILD ADD . /elements
