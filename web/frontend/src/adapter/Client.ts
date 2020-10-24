interface Configuration {
  api_endpoint: string;
}

export interface MessageTemplates {
  templates: string[];
}

export class Client {
  public readonly apiEndpoint: string;

  constructor(api_endpoint: string) {
    this.apiEndpoint = api_endpoint;
  }

  public static async configure(): Promise<Client> {
    const res = await fetch('/.well-known/configuration');
    if (res.status !== 200) {
      throw Error('failed to get configuration');
    }
    const config = (await res.json()) as Configuration;

    return new Client(config.api_endpoint);
  }

  public async messageTemplates(): Promise<MessageTemplates> {
    const res = await fetch(`${this.apiEndpoint}/message/templates`, {
      mode: 'cors',
    });
    if (res.status !== 200) {
      throw Error('failed to get message templates');
    }

    return (await res.json()) as MessageTemplates;
  }

  public async setMessage(doorbell_id: string, message: string): Promise<void> {
    const res = await fetch(`${this.apiEndpoint}/message/set`, {
      method: 'POST',
      mode: 'cors',
      body: JSON.stringify({
        doorbell_id,
        message,
        duration_sec: 60,
      }),
    });
    if (res.status !== 201) {
      throw Error('failed to set message templates');
    }
  }
}
