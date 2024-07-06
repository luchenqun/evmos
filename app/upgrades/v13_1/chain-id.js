import { ethers } from 'ethers';

const main = async () => {
  try {
    const rpc = 'http://127.0.0.1:8545';
    const provider = new ethers.JsonRpcProvider(rpc);
    const privateKey = 'f78a036930ce63791ea6ea20072986d8c3f16a6811f6a2583b0787c45086f769';
    const wallet = new ethers.Wallet(privateKey, provider);

    const network = await provider.getNetwork();
    console.log(`fetch chainId: ${network.chainId}\n`);

    const oldChainId = '9000';
    const newChainId = '168168';
    const gasLimit = '21000';
    const gasPrice = '10000000000';

    const to = '0x1111102dd32160b064f2a512cdef74bfdb6a9f96';
    const value = '1000000000000000000';

    try {
      const tx = await wallet.sendTransaction({
        to,
        value,
        gasPrice,
        gasLimit,
      });
      const receipt = await tx.wait();
      console.log(`automatic chain id ${network.chainId} retrieval transaction successful\n`);
    } catch (e) {
      console.log('automatic chain id retrieval transaction fail');
      console.log(e, '\n');
    }

    try {
      const tx = await wallet.sendTransaction({
        to,
        value,
        gasPrice,
        gasLimit,
        chainId: oldChainId,
      });
      const receipt = await tx.wait();
      console.log(`use chain id ${oldChainId} transaction successful\n`);
    } catch (e) {
      console.log(`use chain id ${oldChainId} transaction fail`);
      console.log(e, '\n');
    }

    try {
      const tx = await wallet.sendTransaction({
        to,
        value,
        gasPrice,
        gasLimit,
        chainId: newChainId,
      });
      const receipt = await tx.wait();
      console.log(`use chain id ${newChainId} transaction successful\n`);
    } catch (e) {
      console.log(`use chain id ${newChainId} transaction fail`);
      console.log(e, '\n');
    }
  } catch (error) {
    console.log('error', error);
  }
};

main();
