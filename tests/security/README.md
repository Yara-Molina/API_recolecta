# Pruebas del borde de autorización

Ejercitan el punto donde vive **toda** la autorización del dominio de rutas: el
árbol que registra `RutaRoutes.Run()`. `api_rutas` no valida tokens en ninguno
de sus endpoints, así que lo que no se detenga aquí no se detiene en ningún
sitio.

No tocan red, base de datos ni contenedores: `api_rutas` se sustituye por un
`httptest.Server` que graba lo que recibe.

## Cómo ejecutarlas

```bash
go test ./tests/security/... ./src/core/... -v
```

Sin Go instalado, desde el contenedor del backend:

```bash
docker exec -it gin_backend_dev go test ./tests/security/... ./src/core/... -v
```

## Cómo leer el resultado

Cada prueba afirma el comportamiento **seguro deseado**, no el actual. Por eso
varias fallan hoy: **una prueba en rojo es un hallazgo, no una prueba rota**.
El mensaje de fallo lleva el identificador (`HALLAZGO F1`, `F2`…), la causa y
qué debería ocurrir en su lugar.

Las pruebas que pasan también importan: documentan defensas que sí existen
(rechazo de `alg=none`, de firma ajena y de token caducado) para que nadie las
retire por descuido.
