import "preact/debug";
import { render } from "preact";
import "@picocss/pico";
import { App } from "./app.jsx";

render(<App />, document.getElementById("app"));
