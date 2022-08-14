/* @refresh reload */
import { render } from "solid-js/web"

import "./index.css"
import { HopeProvider } from "@hope-ui/solid"
import App from "./App"

const config = {
  initialColorMode: "system",
}

render(
  () => (
    <HopeProvider config={config}>
      <App />
    </HopeProvider>
  ),
  document.getElementById("root")
)
