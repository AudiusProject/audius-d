import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Layout from "./Layout";
import Status from "../pages/Status";
import Network from "../pages/Network";
import { TxViewer } from "../pages/TX";
import { Ddex } from "../pages/ddex";

const App: React.FC = () => {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Status />} />
          <Route path="/network" element={<Network />} />
          <Route path="/explorer" element={<TxViewer />} />
          <Route path="/ddex" element={<Ddex />} />
        </Routes>
      </Layout>
    </Router>
  );
};

export default App;
