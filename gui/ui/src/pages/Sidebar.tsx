import React from "react";
import { Link } from "react-router-dom";

const Sidebar: React.FC = () => {
  return (
    <aside className="w-1/4 p-4">
      <nav>
        <ul>
          <li>
            <Link to="/">Status</Link>
          </li>
          <li>
            <Link to="/network">Network</Link>
          </li>
        </ul>
      </nav>
    </aside>
  );
};

export default Sidebar;
