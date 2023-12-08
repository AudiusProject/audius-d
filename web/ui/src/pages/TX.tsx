import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Link, useLocation, useSearchParams } from "react-router-dom";
import {
  EM_ADDRESS,
  decodeEmLog,
  useEthersProvider,
  useSomeDiscoveryEndpoint,
} from "../utils/acdc-client";
import { useEnvVars } from "../providers/EnvVarsProvider";
import { TxDetail } from "./TXDetail";

interface Log {
  address: string;
  blockHash: string;
  blockNumber: number;
  data: string;
  provider: any;
  transactionHash: string;
}

const TxRow = ({
  log,
  idx,
  discoveryEndpoint,
  onSelect,
  isSelected,
}: {
  log: Log;
  idx: number;
  discoveryEndpoint: string;
  onSelect: (txHash: string) => void;
  isSelected: boolean;
}) => {
  const rowStyle = isSelected ? "bg-sky-100" : "";
  const em = decodeEmLog(log.data);
  return (
    <tr
      key={idx}
      className={`hover:bg-sky-50 ${rowStyle}`}
      onClick={() => {
        console.log(log, em);
      }}
    >
      <td className="tableCell">
        <div className="hover:underline">
          <Link to={`?block=${log.blockNumber.toString()}`}>
            {log.blockNumber.toString()}
          </Link>
        </div>
      </td>
      <td className="tableCellFirst">
        <div
          className="cursor-pointer text-blue-600 hover:text-blue-800 underline"
          onClick={() => onSelect(log.transactionHash)}
        >
          {log.transactionHash.substring(0, 15)}...
        </div>
      </td>
      <td className="tableCell">
        <div className="text-xs text-blue-600 hover:text-blue-800 underline">
          <a
            href={`${discoveryEndpoint}/users/account?wallet=${em._signer}`}
            target="_blank"
            rel="noreferrer"
          >
            {em._signer}
          </a>
        </div>
      </td>
      <td className="tableCell">
        <div className="text-blue-600 hover:text-blue-800 underline">
          <a
            href={`${discoveryEndpoint}/users?id=${(
              em._userId as number
            ).toString()}`}
            target="_blank"
            rel="noreferrer"
          >
            {(em._userId as number).toString()}
          </a>
        </div>
      </td>
      <td className="tableCell">{em._action}</td>
      <td className="tableCell">{em._entityType}</td>
      <td className="tableCell">{(em._entityId as number).toString()}</td>
      <td className="tableCell">
        <pre className="text-xs">{em._metadata}</pre>
      </td>
    </tr>
  );
};

export function TxViewer() {
  const location = useLocation();
  const [searchParams, setSearchParams] = useSearchParams();
  const discoveryEndpoint = useSomeDiscoveryEndpoint();
  const provider = useEthersProvider();
  const { env } = useEnvVars();
  const isProd = env === "prod";
  const [selectedTx, setSelectedTx] = useState("");
  const handleSelect = (txHash: string) => {
    setSelectedTx(txHash);
  };
  const handleCloseModal = () => {
    setSelectedTx("");
  };

  const { data, isLoading } = useQuery({
    queryKey: [isProd, location.pathname, location.search],
    queryFn: async () => {
      let latestBlock = parseInt(searchParams.get("block") || "");
      if (!latestBlock) latestBlock = await provider.getBlockNumber();

      const logs: Log[] = await provider.getLogs({
        fromBlock: latestBlock - (isProd ? 1000 : 10000),
        toBlock: latestBlock,
        address: EM_ADDRESS,
      });

      logs.reverse();
      return { latestBlock, logs };
    },
  });

  if (isLoading || !data) return <div>loading</div>;
  const { latestBlock, logs } = data;

  function showOlder() {
    const no = logs[logs.length - 1].blockNumber;
    searchParams.set("block", no.toString());
    setSearchParams(searchParams);
  }

  function setBlock(b: string) {
    searchParams.set("block", b);
    setSearchParams(searchParams);
  }

  return (
    <div className="flex flex-col space-y-4">
      <div className="font-bold text-xl">Recent Transactions</div>

      <div className="my-2 flex items-center gap-2">
        <div>Block:</div>
        <input
          onChange={(e) => setBlock(e.target.value)}
          placeholder="block number"
          className="p-2 my-2 rounded niceBorder"
          value={latestBlock.toString()}
        />
      </div>

      <div className="mt-8 flow-root">
        <div className="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
            <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
              <table className="min-w-full divide-y divide-gray-300">
                <thead className="bg-gray-50">
                  <tr>
                    <th scope="col" className="tableHeaderCellFirst">
                      block no
                    </th>
                    <th scope="col" className="tableHeaderCell">
                      tx hash
                    </th>
                    <th scope="col" className="tableHeaderCell">
                      signed by
                    </th>
                    <th scope="col" className="tableHeaderCell">
                      user id
                    </th>
                    <th scope="col" className="tableHeaderCell">
                      action
                    </th>
                    <th scope="col" className="tableHeaderCell">
                      type
                    </th>
                    <th scope="col" className="tableHeaderCell">
                      id
                    </th>
                    <th scope="col" className="tableHeaderCell">
                      metadata
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 bg-white">
                  {logs.map((log, idx) => (
                    <TxRow
                      key={idx}
                      idx={idx}
                      log={log}
                      discoveryEndpoint={discoveryEndpoint}
                      onSelect={handleSelect}
                      isSelected={selectedTx === log.transactionHash}
                    />
                  ))}
                </tbody>
              </table>
              <div className="w-1/2">
                {selectedTx && (
                  <TxDetail tx={selectedTx} onClose={handleCloseModal} />
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      <button className="btn btn-blue w-24" onClick={showOlder}>
        Older
      </button>
    </div>
  );
}
