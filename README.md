# Luda.Farm Libs

## Mailer

A simple interface for sending email over smtp with Login Auth.
Automatically BCCs all recipients.
Content type is set to HTML, UTF-8.

## Router

Simple HTTP router that implements the `net/http#Handler` interface.
Uses `$` prefix for URL parameters.
Automatic CORS handling with configurable allowed origins.
