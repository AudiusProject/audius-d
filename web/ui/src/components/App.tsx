import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Layout from "./Layout";
import Status from "../pages/Status";
import Network from "../pages/Network";

const App: React.FC = () => {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Status />} />
          <Route path="/network" element={<Network />} />
        </Routes>
      </Layout>
    </Router>
  );
};

export default App;
