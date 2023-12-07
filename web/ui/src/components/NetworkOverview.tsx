import { useEffect, useState } from "react";
import BN from "bn.js";
import useSWR from "swr";
import { Tracker, Color } from "@tremor/react";
import { useEnvVars } from "../providers/EnvVarsProvider";
import { useAudiusLibs } from "../providers/AudiusLibsProvider";
import type { AudiusLibs } from "@audius/sdk/dist/WebAudiusLibs.d.ts";
import { formatWei } from "../utils/helpers";
import useNodes from "../hooks/useNodes";
import { UptimeResponse } from "../components/Uptime";

const DEPLOYER_CUT_BASE = new BN("100");

interface NodeResponse {
  blockNumber: number;
  delegateOwnerWallet: string;
  endpoint: string;
  owner: string;
  spID: number;
  type: string;
}

interface HealthResponse {
  data: {
    version: string;
    discovery_provider_healthy?: boolean;
    healthy?: boolean;
  };
}

interface DiscoveryRequestCountResponse {
  timestamp: string;
  unique_count: number;
  total_count: number;
}

interface ContentRequestCountResponse {
  timestamp: string;
  count: number;
}

interface RequestCountsResponse {
  data: DiscoveryRequestCountResponse[] | ContentRequestCountResponse[];
}

interface Tracker {
  color: Color;
  tooltip: string;
}

type BigNumber = BN;

type ServiceProvider = {
  deployerCut: number;
  deployerStake: BigNumber;
  maxAccountStake: BigNumber;
  minAccountStake: BigNumber;
  numberOfEndpoints: number;
  validBounds: boolean;
};

type GetPendingDecreaseStakeRequestResponse = {
  lockupExpiryBlock: number;
  amount: BN;
};

type GetPendingUndelegateRequestResponse = {
  amount: BigNumber;
  lockupExpiryBlock: number;
  target: string;
};

type GetIncreaseDelegateStakeEventsResponse = {
  blockNumber: number;
  delegator: string;
  increaseAmount: BN;
  serviceProvider: string;
};

type Delegate = {
  wallet: string;
  amount: BN;
  activeAmount: BN;
  name?: string;
  img?: string;
};

type User = {
  wallet: string;
  delegates: Array<Delegate>;
  totalDelegatorStake: BigNumber | undefined;
  pendingUndelegateRequest: GetPendingUndelegateRequestResponse | undefined;
};

type Operator = {
  serviceProvider: ServiceProvider;
  delegators: Array<Delegate>;
  pendingDecreaseStakeRequest: GetPendingDecreaseStakeRequestResponse;
} & User;

class FetchError extends Error {
  info: any;
  status: number;

  constructor(message: string) {
    super(message);
    this.info = null;
    this.status = 0;
  }
}

const fetcher = async (url: string) => {
  const res = await fetch(url);

  // If the status code is not in the range 200-299,
  // we still try to parse and throw it.
  if (!res.ok) {
    const error = new FetchError("An error occurred while fetching the data.");
    // Attach extra info to the error object.
    error.info = await res.json();
    error.status = res.status;
    throw error;
  }

  return res.json();
};

const UptimeTracker = ({ data }: { data: UptimeResponse }) => {
  if (!data?.uptime_raw_data) {
    return null;
  }

  const trackerData: Tracker[] = [];
  for (const [hour, up] of Object.entries(data?.uptime_raw_data)) {
    const hourString = new Date(hour).toUTCString();
    if (up) {
      trackerData.push({
        color: "emerald",
        tooltip: `${hourString}: Operational`,
      });
    } else {
      trackerData.push({
        color: "rose",
        tooltip: `${hourString}: Down`,
      });
    }
  }

  return <Tracker data={trackerData} className="w-20 mx-auto mt-2" />;
};

// Operator metadata helpers

const getUserDelegates = async (delegator: string, audiusLibs: AudiusLibs) => {
  const delegates = [];
  const increaseDelegateStakeEvents =
    await audiusLibs.ethContracts?.DelegateManagerClient.getIncreaseDelegateStakeEvents(
      {
        delegator,
        serviceProvider: "",
        queryStartBlock: 0,
      },
    );
  const pendingUndelegateRequest =
    await audiusLibs.ethContracts?.DelegateManagerClient.getPendingUndelegateRequest(
      delegator,
    );
  let serviceProviders = (
    increaseDelegateStakeEvents as GetIncreaseDelegateStakeEventsResponse[]
  ).map((e) => e.serviceProvider);
  serviceProviders = [...new Set(serviceProviders)];
  for (const sp of serviceProviders) {
    const delegators =
      await audiusLibs.ethContracts?.DelegateManagerClient.getDelegatorsList(
        sp,
      );
    if ((delegators as Array<string>).includes(delegator)) {
      const amountDelegated =
        await audiusLibs.ethContracts?.DelegateManagerClient.getDelegatorStakeForServiceProvider(
          delegator,
          sp,
        );
      let activeAmount = amountDelegated;

      if (
        pendingUndelegateRequest!.lockupExpiryBlock !== 0 &&
        pendingUndelegateRequest!.target === sp
      ) {
        activeAmount = activeAmount!.sub(pendingUndelegateRequest!.amount);
      }

      delegates.push({
        wallet: sp,
        amount: amountDelegated,
        activeAmount,
      });
    }
  }
  return [...delegates];
};

const getDelegatorAmounts = async (
  wallet: string,
  audiusLibs: AudiusLibs,
): Promise<
  Array<{
    wallet: string;
    amount: BN;
    activeAmount: BN;
    // name?: string
    // img: string
  }>
> => {
  const delegators =
    await audiusLibs.ethContracts?.DelegateManagerClient.getDelegatorsList(
      wallet,
    );
  const delegatorAmounts = [];
  for (const delegatorWallet of delegators) {
    const amountDelegated =
      await audiusLibs.ethContracts?.DelegateManagerClient.getDelegatorStakeForServiceProvider(
        delegatorWallet as string,
        wallet,
      );
    let activeAmount = amountDelegated;
    const pendingUndelegateRequest =
      await audiusLibs.ethContracts?.DelegateManagerClient.getPendingUndelegateRequest(
        delegatorWallet as string,
      );

    if (
      pendingUndelegateRequest!.lockupExpiryBlock !== 0 &&
      pendingUndelegateRequest!.target === wallet
    ) {
      activeAmount = activeAmount!.sub(pendingUndelegateRequest!.amount);
    }

    delegatorAmounts.push({
      wallet: delegatorWallet,
      amount: amountDelegated!,
      activeAmount: activeAmount!,
    });
  }
  return delegatorAmounts;
};

const getUserMetadata = async (
  wallet: string,
  audiusLibs: AudiusLibs,
): Promise<User> => {
  const delegates = await getUserDelegates(wallet, audiusLibs);
  const totalDelegatorStake =
    await audiusLibs.ethContracts?.DelegateManagerClient.getTotalDelegatorStake(
      wallet,
    );
  const pendingUndelegateRequest =
    await audiusLibs.ethContracts?.DelegateManagerClient.getPendingUndelegateRequest(
      wallet,
    );

  const user = {
    wallet,
    delegates: delegates as Delegate[],
    totalDelegatorStake,
    pendingUndelegateRequest,
  };

  return user;
};

const getServiceProviderMetadata = async (
  wallet: string,
  audiusLibs: AudiusLibs,
) => {
  const totalStakedFor =
    await audiusLibs.ethContracts?.StakingProxyClient.totalStakedFor(wallet);
  const delegatedTotal =
    await audiusLibs.ethContracts?.DelegateManagerClient.getTotalDelegatedToServiceProvider(
      wallet,
    );
  const delegators = await getDelegatorAmounts(wallet, audiusLibs);
  delegators.sort((a, b) => (b.activeAmount.gt(a.activeAmount) ? 1 : -1));

  const serviceProvider: ServiceProvider | undefined =
    await audiusLibs.ethContracts?.ServiceProviderFactoryClient.getServiceProviderDetails(
      wallet,
    );
  const pendingDecreaseStakeRequest =
    await audiusLibs.ethContracts?.ServiceProviderFactoryClient.getPendingDecreaseStakeRequest(
      wallet,
    );

  return {
    serviceProvider,
    pendingDecreaseStakeRequest,
    totalStakedFor,
    delegatedTotal,
    delegators,
  };
};

const getOperatorMetadata = async (wallet: string, audiusLibs: AudiusLibs) => {
  const user = await getUserMetadata(wallet, audiusLibs);
  const serviceProvider = await getServiceProviderMetadata(wallet, audiusLibs);
  const operator = {
    ...user,
    ...serviceProvider,
  };

  return operator;
};

// Rewards helpers

/**
 * Calculates the net minted amount for a service operator prior to
 * distribution among the service provider and their delegators.
 * Reference processClaim in the claims manager contract.
 * NOTE: minted amount is calculated using values at the init claim block
 *
 * wallet The service operator's wallet address
 * totalLocked The total token currently locked (decrease stake and delegation)
 * blockNumber The blocknumber of the claim to process
 * fundingAmount The amount of total funds allocated per claim round
 * The net minted amount
 */
const getMintedAmountAtBlock = async ({
  wallet,
  totalLocked,
  blockNumber,
  fundingAmount,
  audiusLibs,
}: {
  wallet: string;
  totalLocked: BN;
  blockNumber: number;
  fundingAmount: BN;
  audiusLibs: AudiusLibs;
}) => {
  const totalStakedAtFundBlockForClaimer =
    await audiusLibs.ethContracts?.StakingProxyClient?.totalStakedForAt(
      wallet,
      blockNumber.toString(),
    );
  const totalStakedAtFundBlock =
    await audiusLibs.ethContracts?.StakingProxyClient?.totalStakedAt(
      blockNumber,
    );
  const activeStake = totalStakedAtFundBlockForClaimer!.sub(totalLocked);
  const rewardsForClaimer = activeStake
    .mul(fundingAmount)
    .div(totalStakedAtFundBlock!);

  return rewardsForClaimer;
};

// Get the operator's active stake = total staked - pending decrease stake + total delegated to operator - operator's delegators' pending decrease stake
const getOperatorTotalActiveStake = (user: Operator) => {
  const userActiveStake = user.serviceProvider.deployerStake.sub(
    user.pendingDecreaseStakeRequest?.amount ?? new BN("0"),
  );
  const userActiveDelegated = user.delegators.reduce((total, delegator) => {
    return total.add(delegator.activeAmount);
  }, new BN("0"));
  const totalActiveStake = userActiveStake.add(userActiveDelegated);
  return totalActiveStake;
};

// Get the amount locked - pending decrease stake, and the operator's delegator's pending decrease delegation
export const getOperatorTotalLocked = (user: Operator) => {
  const lockedPendingDecrease =
    user.pendingDecreaseStakeRequest?.amount ?? new BN("0");
  const lockedDelegation = user.delegators.reduce((totalLocked, delegate) => {
    return totalLocked.add(delegate.amount.sub(delegate.activeAmount));
  }, new BN("0"));
  const totalLocked = lockedPendingDecrease.add(lockedDelegation);
  return totalLocked;
};

const getOperatorRewards = ({
  user,
  totalRewards,
  deployerCutBase = DEPLOYER_CUT_BASE,
}: {
  user: Operator;
  totalRewards: BN;
  deployerCutBase?: BN;
}) => {
  const totalActiveStake = getOperatorTotalActiveStake(user);
  const deployerCut = new BN(user.serviceProvider.deployerCut);

  const totalDelegatedRewards = user.delegators.reduce((total, delegate) => {
    const delegateRewards = getDelegateRewards({
      delegateAmount: delegate.activeAmount,
      totalRoundRewards: totalRewards,
      totalActive: totalActiveStake,
      deployerCut,
      deployerCutBase,
    });
    return total.add(delegateRewards.delegatorCut);
  }, new BN("0"));

  const operatorRewards = totalRewards.sub(totalDelegatedRewards);
  return operatorRewards;
};

const getDelegateRewards = ({
  delegateAmount,
  totalRoundRewards,
  totalActive,
  deployerCut,
  deployerCutBase = DEPLOYER_CUT_BASE,
}: {
  delegateAmount: BN;
  totalRoundRewards: BN;
  totalActive: BN;
  deployerCut: BN;
  deployerCutBase?: BN;
}) => {
  const rewardsPriorToSPCut = delegateAmount
    .mul(totalRoundRewards)
    .div(totalActive);
  const spDeployerCut = delegateAmount
    .mul(totalRoundRewards)
    .mul(deployerCut)
    .div(totalActive)
    .div(deployerCutBase);
  return {
    spCut: spDeployerCut,
    delegatorCut: rewardsPriorToSPCut.sub(spDeployerCut),
  };
};

/**
 * Calculates and returns the total rewards for a user from a
 * claim given a blocknumber
 *
 * fundsPerRound The amount of rewards given out in a round
 * blockNumber The block number to process the claim event for
 * expected rewards for the user at the claim block
 */
const getRewardForClaimBlock = async ({
  user,
  fundsPerRound,
  blockNumber,
  audiusLibs,
}: {
  user: User | Operator;
  fundsPerRound: BN;
  blockNumber: number;
  audiusLibs: AudiusLibs;
}): Promise<BN> => {
  let totalRewards = new BN("0");

  // If the user is a service provider, retrieve their expected rewards for staking
  if ("serviceProvider" in user) {
    const lockedPendingDecrease =
      user.pendingDecreaseStakeRequest?.amount ?? new BN("0");
    const lockedDelegation =
      await audiusLibs.ethContracts?.DelegateManagerClient.getTotalLockedDelegationForServiceProvider(
        user.wallet,
      );
    const totalLocked = lockedPendingDecrease.add(lockedDelegation!);
    const mintedRewards = await getMintedAmountAtBlock({
      totalLocked,
      fundingAmount: fundsPerRound,
      wallet: user.wallet,
      blockNumber,
      audiusLibs,
    });
    const operatorRewards = getOperatorRewards({
      user: user,
      totalRewards: mintedRewards,
    });
    totalRewards = totalRewards.add(operatorRewards);
  }

  // For each service operator the user delegates to, calculate the expected rewards for delegating
  for (const delegate of (user as User).delegates) {
    const operator = await getOperatorMetadata(delegate.wallet, audiusLibs);
    const deployerCut = new BN(
      operator.serviceProvider!.deployerCut.toString(),
    );
    const operatorActiveStake = getOperatorTotalActiveStake(
      operator as Operator,
    );
    const operatorTotalLocked = getOperatorTotalLocked(operator as Operator);
    const userMintedRewards = await getMintedAmountAtBlock({
      totalLocked: operatorTotalLocked,
      fundingAmount: fundsPerRound,
      wallet: delegate.wallet,
      blockNumber,
      audiusLibs,
    });
    const delegateRewards = getDelegateRewards({
      delegateAmount: delegate.activeAmount,
      totalRoundRewards: userMintedRewards,
      totalActive: operatorActiveStake,
      deployerCut,
    });
    totalRewards = totalRewards.add(delegateRewards.delegatorCut);
  }
  return totalRewards;
};

/**
 * Calculates and returns the active stake for address
 *
 * Active stake = (active deployer stake + active delegator stake)
 *      active deployer stake = (direct deployer stake - locked deployer stake)
 *          locked deployer stake = amount of pending decreaseStakeRequest for address
 *      active delegator stake = (total delegator stake - locked delegator stake)
 *          locked delegator stake = amount of pending undelegateRequest for address
 */
const calculateActiveStake = (operator: Operator): BN => {
  let activeDeployerStake = new BN("0");
  let activeDelegatorStake = new BN("0");
  if ("serviceProvider" in operator) {
    const deployerStake = operator.serviceProvider?.deployerStake;
    const {
      amount: pendingDecreaseStakeAmount = new BN("0"),
      lockupExpiryBlock = 0,
    } = operator.pendingDecreaseStakeRequest ?? {};

    if (deployerStake) {
      if (lockupExpiryBlock !== 0) {
        activeDeployerStake = deployerStake.sub(pendingDecreaseStakeAmount);
      } else {
        activeDeployerStake = deployerStake;
      }
    }
  }

  if (operator.pendingUndelegateRequest?.lockupExpiryBlock !== 0) {
    // Ensure operator.totalDelegatorStake and operator.pendingUndelegateRequest.amount are defined and are BN
    if (
      operator.totalDelegatorStake &&
      operator.pendingUndelegateRequest?.amount
    ) {
      activeDelegatorStake = operator.totalDelegatorStake.sub(
        operator.pendingUndelegateRequest.amount,
      );
    } else {
      activeDelegatorStake = new BN("0");
    }
  } else {
    // If operator.totalDelegatorStake is not defined, default to BN("0")
    activeDelegatorStake = operator.totalDelegatorStake || new BN("0");
  }
  const activeStake = activeDelegatorStake.add(activeDeployerStake);
  return activeStake;
};

const NodeRow = ({
  node,
  nodeType,
  prevDate,
}: {
  node: NodeResponse;
  nodeType: string;
  prevDate: string;
}) => {
  const { audiusLibs } = useAudiusLibs();

  const [bondedData, setBondedData] = useState("");
  const [bondedDataError, setBondedDataError] = useState(false);

  const [rewardData, setRewardData] = useState("");
  const [rewardDataError, setRewardDataError] = useState(false);

  useEffect(() => {
    const fetchAudData = async (wallet: string, audiusLibs: AudiusLibs) => {
      const operator = await getOperatorMetadata(wallet, audiusLibs);
      try {
        const activeStake = calculateActiveStake(operator as Operator);
        setBondedData(formatWei(activeStake));
      } catch (error) {
        console.error(
          `Could not fetch bonded $AUDIO for ${operator.wallet}`,
          error,
        );
        setBondedDataError(true);
      }
      try {
        const blockNumber =
          await audiusLibs.ethWeb3Manager?.web3.eth.getBlockNumber();
        const fundsPerRound =
          await audiusLibs.ethContracts?.ClaimsManagerClient.getFundsPerRound();
        const weeklyReward = await getRewardForClaimBlock({
          user: operator,
          fundsPerRound: new BN(fundsPerRound!),
          blockNumber: blockNumber!,
          audiusLibs,
        });
        const est24hReward = weeklyReward.div(new BN(7));
        setRewardData(formatWei(est24hReward));
      } catch (error) {
        console.error(
          `Could not fetch weekly $AUDIO reward for ${operator.wallet}`,
          error,
        );
        setRewardDataError(true);
      }
    };

    if (node.owner && audiusLibs) {
      void fetchAudData(node.owner, audiusLibs);
    }
  }, [node, audiusLibs]);

  // fetch health
  const { data: healthData, error: healthDataError } = useSWR(
    `${node.endpoint}/health_check?enforce_block_diff=true&healthy_block_diff=250&plays_count_max_drift=720`,
    fetcher,
  ) as { data: HealthResponse; error: any };
  const health = healthData?.data;

  // fetch 24h uptime data
  const { data: uptimeData, error: uptimeDataError } = useSWR(
    `${node.endpoint}/d_api/uptime?host=${node.endpoint}&durationHours=12`,
    fetcher,
  ) as { data: UptimeResponse; error: any };

  // fetch request counts from previous day
  let requestCount;
  const requestsPath =
    nodeType == "discovery"
      ? "v1/metrics/routes/week?bucket_size=day"
      : "internal/metrics/blobs-served/week?bucket_size=day";
  const { data: requestsData, error: requestsDataError } = useSWR(
    `${node.endpoint}/${requestsPath}`,
    fetcher,
  ) as { data: RequestCountsResponse; error: any };
  const requests = requestsData?.data;
  if (requests && requests.length > 0) {
    if (nodeType == "discovery") {
      const lastDay = (requests as DiscoveryRequestCountResponse[])[
        requests.length - 1
      ];
      if (lastDay.timestamp != prevDate) {
        requestCount = 0;
      } else {
        requestCount = lastDay.total_count;
      }
    } else {
      const lastDay = (requests as ContentRequestCountResponse[])[
        requests.length - 1
      ];
      if (lastDay.timestamp.replace("T00:00:00Z", "") != prevDate) {
        requestCount = 0;
      } else {
        requestCount = (requests as ContentRequestCountResponse[])[
          requests.length - 1
        ].count;
      }
    }
  }

  return (
    <tr>
      <td className="tableCellFirst">
        <div className="flex items-center justify-center">
          {!healthDataError && !health ? (
            "loading..."
          ) : healthDataError ? (
            <span className="flex w-3 h-3 me-3 bg-red-500 rounded-full"></span>
          ) : health?.healthy || health?.discovery_provider_healthy ? (
            <span className="flex w-3 h-3 me-3 bg-green-500 rounded-full"></span>
          ) : (
            <span className="flex w-3 h-3 me-3 bg-red-500 rounded-full"></span>
          )}
        </div>
      </td>
      <td className="tableCell">
        {!healthDataError && !health
          ? "loading..."
          : healthDataError
            ? "error"
            : health?.version}
      </td>
      <td className="tableCell">
        {!uptimeDataError && !uptimeData ? (
          "loading..."
        ) : uptimeDataError || (uptimeData && !uptimeData.uptime_raw_data) ? (
          "error"
        ) : (
          <UptimeTracker key={node.endpoint} data={uptimeData} />
        )}
      </td>
      <td className="tableCell">{node.endpoint}</td>
      <td className="tableCell">
        {!bondedDataError && !bondedData
          ? "loading..."
          : bondedDataError
            ? "error"
            : bondedData}
      </td>
      <td className="tableCell">
        {!rewardDataError && !rewardData
          ? "loading..."
          : rewardDataError
            ? "error"
            : rewardData}
      </td>
      <td className="tableCell">
        {!requestsDataError && requestCount == undefined
          ? "loading..."
          : requestsDataError
            ? "error"
            : requestCount}
      </td>
      <td className="tableCell">{node.owner}</td>
    </tr>
  );
};

const NetworkOverview = () => {
  const { nodeType } = useEnvVars();
  const {
    data: nodes,
    isPending: isListNodesPending,
    error: listNodesError,
  } = useNodes(nodeType);

  // For requests header
  const prevDate = new Date();
  prevDate.setDate(prevDate.getDate() - 1);
  const prevDateString = prevDate.toISOString().substring(0, 10);

  return (
    <>
      {isListNodesPending ? (
        "loading..."
      ) : listNodesError ? (
        "error"
      ) : (
        <div className="mt-8 flow-root">
          <div className="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
            <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
              <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
                <table className="min-w-full divide-y divide-gray-300">
                  <thead className="bg-gray-50">
                    <tr>
                      <th scope="col" className="tableHeaderCellFirst">
                        Health
                      </th>
                      <th scope="col" className="tableHeaderCell">
                        Version
                      </th>
                      <th scope="col" className="tableHeaderCell">
                        Uptime
                      </th>
                      <th scope="col" className="tableHeaderCell">
                        Host
                      </th>
                      <th scope="col" className="tableHeaderCell">
                        Bond $AUDIO
                      </th>
                      <th scope="col" className="tableHeaderCell">
                        Reward (24H)
                      </th>
                      <th scope="col" className="tableHeaderCell">
                        Requests ({prevDateString})
                      </th>
                      <th scope="col" className="tableHeaderCell">
                        Operator
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200 bg-white">
                    {(nodes as NodeResponse[]).map((node) => (
                      <NodeRow
                        key={node.endpoint}
                        node={node}
                        nodeType={nodeType}
                        prevDate={prevDateString}
                      />
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default NetworkOverview;
