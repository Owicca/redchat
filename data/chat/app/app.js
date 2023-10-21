import React, {useEffect, useState} from "react";

export function App({getClient}) {
    const [client, setClient] = useState(null);
    const [data, setData] = useState("");

    useEffect(() => {
        (async (getClient) => {
            let cl = await getClient("10.247.12.224", "6379", "user", "user");
            setClient(cl);
        })(getClient);
    }, [getClient]);

    useEffect(() => {
        (async (client) => {
            if(client == null) {
                return;
            }

            let d = await client.get("test1");
            setData(d);
        })(client);
    }, [client]);

    console.log("[INFO]", data);

    return (
        <h1>{data}</h1>
    );
}