import { Route, Routes } from "react-router";
import "./App.css";
import Home from "./components/pages/Home";
import { Layout } from "./components/pages/Layout";

function App() {
    return (
        <Routes>
            <Route element={<Layout />}>
                <Route path="/" element={<Home />} />
            </Route>
        </Routes>
    );
}

export default App;
