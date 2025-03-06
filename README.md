## System Documentation (EN)

### Overview
This system receives MP4 video uploads into a Google Cloud Storage (GCP) bucket and converts them into the MPEG-DASH format in a highly available and scalable way using Golang.

### Conversion Process
1. The client uploads a video to Google Cloud Storage.
2. The app receives a message via RabbitMQ specifying which video needs conversion.
3. The app downloads the video.
4. The app fragments the video.
5. The video is converted to MPEG-DASH format.
6. The converted video is uploaded to Google Cloud Storage.
7. A notification is sent to the queue indicating success or failure.
8. In case of failure, the message is forwarded to the Dead Letter Exchange (DLE).

### Data Structure

#### Queue Input
```json
{
  "resource_id": "uuid",
  "path": "path-google-cloud-storage"
}
```

#### Success Output
```json
{
  "id": "uuid-do-job",
  "output_bucket_path": "path-do-bucket",
  "status": "completed",
  "video": {
    "encoded_video_folder": "path-do-video-mpeg-dash",
    "resource_id": "uuid",
    "path": "path-google-cloud-storage"
  },
  "error": "",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

#### Error Output
```json
{
  "message": {
    "resource_id": "uuid",
    "path": "path-google-cloud-storage"
  },
  "error": "error-message"
}
```

### Key Features
- The system processes multiple uploads in parallel and concurrently.
- Each stage of the process is highly scalable.
- The encoded video is composed of multiple small chunks that form the final video.

### Architecture
- **Layer 1 - Domain**: System core, entities, and business logic.
- **Layer 2 - Application**: Uses the Domain layer and executes application logic.
- **Layer 3 - Framework**: Utilities and third-party tools.
```


---

[PT]
# Documentação do Sistema de Processamento de Vídeos

```markdown
## Visão Geral (PT-BR)
Este sistema recebe uploads de vídeos MP4 em um bucket do Google Cloud Storage (GCP) e os converte para o formato MPEG-DASH de forma altamente disponível e escalável, utilizando Golang.

### Processo de Conversão
1. O cliente faz upload de um vídeo no Google Cloud Storage.
2. O aplicativo recebe uma mensagem via RabbitMQ informando qual vídeo deve ser convertido.
3. O aplicativo faz o download do vídeo.
4. O aplicativo fragmenta o vídeo.
5. O vídeo é convertido para o formato MPEG-DASH.
6. O vídeo convertido é enviado para o storage do Google.
7. Uma notificação é enviada na fila indicando sucesso ou erro.
8. Em caso de falha, a mensagem é encaminhada para a Dead Letter Exchange (DLE).

### Estrutura de Dados

#### Entrada da Fila
```json
{
  "resource_id": "uuid",
  "path": "path-google-cloud-storage"
}
```

#### Saída em caso de sucesso
```json
{
  "id": "uuid-do-job",
  "output_bucket_path": "path-do-bucket",
  "status": "completed",
  "video": {
    "encoded_video_folder": "path-do-video-mpeg-dash",
    "resource_id": "uuid",
    "path": "path-google-cloud-storage"
  },
  "error": "",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

#### Saída em caso de erro
```json
{
  "message": {
    "resource_id": "uuid",
    "path": "path-google-cloud-storage"
  },
  "error": "mensagem-do-erro"
}
```

### Características Importantes
- O sistema processa diversos uploads de forma paralela e concorrente.
- Cada etapa do processo é altamente escalável.
- O vídeo convertido é dividido em pequenos segmentos (chunks), formando o vídeo final.

### Arquitetura
- **Layer 1 - Domain**: Core do sistema, entidades e regras de negócio.
- **Layer 2 - Application**: Usa a camada Domain e executa a lógica do aplicativo.
- **Layer 3 - Framework**: Utilitários e ferramentas de terceiros.

