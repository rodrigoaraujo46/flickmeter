import { Route, Routes } from "react-router";
import "./App.css";
import { ErrorBoundary } from "./components/ErrorBoundary";
import { Layout } from "./components/Layout";
import Home from "./pages/Home";
import Movie from "./pages/Movie";

function App() {
    return (
        <ErrorBoundary fallback={<div>Error happened</div>}>
            <Routes>
                <Route element={<Layout />}>
                    <Route index element={<Home />} />
                    <Route path="/movies/:id" element={<Movie />} />
                </Route>
            </Routes>
        </ErrorBoundary>
    );
}

export default App;
