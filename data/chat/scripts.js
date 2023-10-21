import {createRoot} from 'react-dom/client';
import {App} from "./app/app";
import {getClient} from "./app/models/redis";

const container = document.getElementById("root");
const root = createRoot(container);

root.render(<App getClient={getClient}/>);