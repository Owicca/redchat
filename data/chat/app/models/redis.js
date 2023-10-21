import {createClient} from 'redis';

export async function getClient(host, port, user, pass) {
    const client = createClient();

    client.on('error', (err) => {
        console.log("[ERROR] herer ", err);
    });

    await client.connect();

    console.log("[INFO]", client);

    return client;
}