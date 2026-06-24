# label-tui-sb1

Terminal UI para imprimir etiquetas desde **SAP Business One** a impresoras **Zebra**.

Conecta con SAP B1 Service Layer, busca artículos por código o descripción, selecciona una plantilla ZPL y envía las etiquetas directamente a la impresora por puerto serie/USB.

## Screens

- **Welcome** → **Login** (credenciales SAP) → **Search** (buscar/seleccionar artículos) → **Select Template** (elegir `.zpl`) → **Preview & Print**

## Features

- Inicio de sesión contra SAP B1 Service Layer (REST API)
- Búsqueda de artículos por código o descripción (`$filter` OData)
- Selección de cantidad por artículo
- Listado de plantillas ZPL desde `~/.label-tui/templates/`
- Renderizado de variables: `{{CODE}}`, `{{DESCRIPTION}}`, `{{BARCODE}}`, `{{PRICE}}`, `{{QTY}}`
- Impresión directa a Zebra por puerto serie (9600 baud, 8N1)
- Persistencia de configuración en `~/.label-tui/settings.json`

## Installation

```bash
go install github.com/JohnDevRD/label-tui-sb1/cmd/label-tui-sb1@latest
```

O build local:

```bash
git clone https://github.com/JohnDevRD/label-tui-sb1.git
cd label-tui-sb1
go build -o label-tui ./cmd/label-tui-sb1
```

## Usage

```bash
./label-tui
```

Configuración persistente (se crea tras el primer uso):

```json
~/.label-tui/settings.json
{
  "company_db": "SBODemoCL",
  "sap_service_layer_url": "http://your-server:50000/b1s/v1",
  "usb_port": "/dev/ttyUSB0",
  "default_template": "etiqueta.zpl"
}
```

Plantillas ZPL en `~/.label-tui/templates/`:

```zpl
^XA
^FO50,50^ADN,36,20^FD{{CODE}}^FS
^FO50,120^ADN,18,10^FD{{DESCRIPTION}}^FS
^FO50,190^BCN,80,Y,N,N^FD{{BARCODE}}^FS
^FO50,300^ADN,18,10^FDPrecio: ${{PRICE}}^FS
^FO50,350^ADN,18,10^FDCant: {{QTY}}^FS
^XZ
```

## Project Structure

```
├── cmd/label-tui-sb1/    # Entry point
├── internal/
│   ├── core/             # Models, SAP client, settings, templates
│   ├── printers/         # Zebra serial communication
│   └── tui/              # Bubbletea UI components & styles
├── .gitignore
├── CODE_OF_CONDUCT.md
├── CONTRIBUTING.md
└── LICENSE
```

## Requirements

- Go 1.24+
- SAP Business One Service Layer (REST API habilitada)
- Impresora Zebra con puerto USB/serie
- Linux (la librería `go.bug.st/serial` tiene soporte multiplataforma, pero el path del puerto depende del OS)

## License

[MIT](LICENSE)
