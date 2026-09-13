#!/usr/bin/env bash

SOLR_URL="http://localhost:8983/solr/products/select"

# CONSULTAS 3 Y 4 — RELEVANCIA
#
# Usá el texto asignado y escribí dos consultas independientes:
#   1. Una búsqueda eDisMax en title y description, que devuelva id, title y
#      score.
#   2. La misma búsqueda, pero dando mayor peso a title mediante qf.
#
# No agregues el boost como filtro: qf define campos y pesos de búsqueda, no
# condiciones obligatorias. Ejecutá ambas consultas y compará numFound, score
# y especialmente el orden de los resultados en el formulario.

# Sin boost: title y description con el mismo peso.
curl -fsS -G "$SOLR_URL" \
  --data-urlencode 'q=running' \
  --data-urlencode 'defType=edismax' \
  --data-urlencode 'qf=title description' \
  --data-urlencode 'fl=id,title,score' \
  --data-urlencode 'wt=json'
printf '\n'

# Con boost: title pesa más que description.
curl -fsS -G "$SOLR_URL" \
  --data-urlencode 'q=running' \
  --data-urlencode 'defType=edismax' \
  --data-urlencode 'qf=title^2 description' \
  --data-urlencode 'fl=id,title,score' \
  --data-urlencode 'wt=json'
printf '\n'
