import React, { useEffect, useState } from "react";
import { Table } from "antd";

interface Container {
  id: number;
  ip: string;
  ping_time: string;
  last_check: string;
}

const App: React.FC = () => {
  const [containers, setContainers] = useState<Container[]>([]);

  useEffect(() => {
    fetch("/containers")
      .then((res) => res.json())
      .then((data) => setContainers(data));
  }, []);

  const columns = [
    { title: "IP", dataIndex: "ip", key: "ip" },
    { title: "Ping Time", dataIndex: "ping_time", key: "ping_time" },
    { title: "Last Check", dataIndex: "last_check", key: "last_check" },
  ];

  return <Table dataSource={containers} columns={columns} rowKey="id" />;
};

export default App;
