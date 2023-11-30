import React from "react";
import Sidebar from "../pages/Sidebar";

const Layout: React.FC = ({ children }) => {
  return (
    <div className="flex">
      <Sidebar />
      <main className="flex-grow p-4">{children}</main>
    </div>
  );
};

export default Layout;
