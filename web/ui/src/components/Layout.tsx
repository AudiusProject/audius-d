import React from "react";
import Sidebar from "../pages/Sidebar";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="flex">
      <Sidebar />
      <main className="flex-grow p-4">{children}</main>
    </div>
  );
};

export default Layout;
