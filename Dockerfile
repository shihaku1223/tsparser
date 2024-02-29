FROM python:3.10-bullseye

RUN --mount=type=ssh ssh -o StrictHostKeyChecking=no git@github.com
