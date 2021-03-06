# Events Protocol Go
Implementação do [protocolo de eventos](https://github.com/GuiaBolso/events-protocol) base. Este pacote não incluirá nem o transporte nem as ferramentas de log, mas apenas uma abstração

# Problema

Gerar o pacote do [protocolo de eventos](https://github.com/GuiaBolso/events-protocol) que possa ser enviado independente do transporte ou do ecossistema

# Como usar

```go
import (
    events "github.com/guiabolso/events-protocol-go
    guuid "github.com/google/uuid"
)

func main() {
    session := events.GenerateEventSession(<UUIDGenerator>)
    eventTemplate := session.RegisterEvent("uuid:event", "1").WithPayload(payload)

    event := eventTemplate.Prepare()
}
```

# Licença

[Apache 2.0](https://github.com/GuiaBolso/events-protocol-go/blob/master/LICENSE) em Guiabolso (r) 2020