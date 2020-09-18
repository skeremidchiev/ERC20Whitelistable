const { accounts, contract, web3 } = require('@openzeppelin/test-environment');
const { BN, expectEvent, expectRevert } = require('@openzeppelin/test-helpers');

const ERC20Whitelistable = contract.fromArtifact('ERC20Whitelistable');

describe('ERC20Whitelistable', function () {
  const [ admin, authorized, otherAuthorized, other, otherAdmin ] = accounts;

  const DEFAULT_ADMIN_ROLE = '0x0000000000000000000000000000000000000000000000000000000000000000';
  const MINTER_ROLE = web3.utils.soliditySha3('MINTER_ROLE');
  const WHITELISTED_ROLE = web3.utils.soliditySha3('WHITELISTED_ROLE');

  beforeEach(async function () {
    this.ERC20Whitelistable = await ERC20Whitelistable.new({ from: admin });
  });

  describe('default admin', function () {
    it('deployer has default admin role', async function () {
      expect(await this.ERC20Whitelistable.hasRole(DEFAULT_ADMIN_ROLE, admin)).to.equal(true);
    });

    it('deployer has minter role', async function () {
      expect(await this.ERC20Whitelistable.hasRole(MINTER_ROLE, admin)).to.equal(true);
    });

    it('deployer has whitelister role', async function () {
      expect(await this.ERC20Whitelistable.hasRole(WHITELISTED_ROLE, admin)).to.equal(true);
    });
  });

  describe('granting', function () {
    // MINTER_ROLE
    it('non-admin cannot grant role to other accounts', async function () {
      await expectRevert(
        this.ERC20Whitelistable.grantRole(MINTER_ROLE, authorized, { from: other }),
        'AccessControl: sender must be an admin to grant',
      );
    });

    it('admin can grant role to other accounts', async function () {
      const receipt = await this.ERC20Whitelistable.grantRole(MINTER_ROLE, authorized, { from: admin });
      expectEvent(receipt, 'RoleGranted', { account: authorized, role: MINTER_ROLE, sender: admin });

      expect(await this.ERC20Whitelistable.hasRole(MINTER_ROLE, authorized)).to.equal(true);
    });

    it('accounts can be granted a role multiple times', async function () {
      await this.ERC20Whitelistable.grantRole(MINTER_ROLE, authorized, { from: admin });
      const receipt = await this.ERC20Whitelistable.grantRole(MINTER_ROLE, authorized, { from: admin });
      expectEvent.notEmitted(receipt, 'RoleGranted');
    });

    // WHITELISTED_ROLE
    it('non-admin cannot grant role to other accounts', async function () {
      await expectRevert(
        this.ERC20Whitelistable.grantRole(WHITELISTED_ROLE, authorized, { from: other }),
        'AccessControl: sender must be an admin to grant',
      );
    });

    it('admin can grant role to other accounts', async function () {
      const receipt = await this.ERC20Whitelistable.grantRole(WHITELISTED_ROLE, authorized, { from: admin });
      expectEvent(receipt, 'RoleGranted', { account: authorized, role: WHITELISTED_ROLE, sender: admin });

      expect(await this.ERC20Whitelistable.hasRole(WHITELISTED_ROLE, authorized)).to.equal(true);
    });

    it('accounts can be granted a role multiple times', async function () {
      await this.ERC20Whitelistable.grantRole(WHITELISTED_ROLE, authorized, { from: admin });
      const receipt = await this.ERC20Whitelistable.grantRole(WHITELISTED_ROLE, authorized, { from: admin });
      expectEvent.notEmitted(receipt, 'RoleGranted');
    });
  });

  describe('mintingAndTransfering', function () {
    const mintAmount = new BN(100);
    const transferAmount = new BN(50);

    it('non-minter cant mint tokens', async function () {
        await expectRevert(
            this.ERC20Whitelistable.mint(other, mintAmount, { from: authorized }),
            'ERC20Whitelistable: must have minter role to mint',
        );
    });

    it('non-whitelisted cant recieve tokens from mint', async function () {
        await this.ERC20Whitelistable.grantRole(MINTER_ROLE, authorized, { from: admin });

        await expectRevert(
            this.ERC20Whitelistable.mint(other, mintAmount, { from: authorized }),
            'ERC20Whitelistable: must be whitelisted to recieve tokens',
        );
    });

    it('mint tokens', async function () {
        await this.ERC20Whitelistable.grantRole(MINTER_ROLE, authorized, { from: admin });
        await this.ERC20Whitelistable.grantRole(WHITELISTED_ROLE, other, { from: admin });
        
        const receipt = await this.ERC20Whitelistable.mint(other, mintAmount, { from: authorized });
        await expectEvent(receipt, "Transfer", { from: '0x0000000000000000000000000000000000000000', to: other, value: mintAmount });

        var result = await this.ERC20Whitelistable.balanceOf(other, { from: other })
        expect(result.toString()).to.equal(mintAmount.toString());

        result = await this.ERC20Whitelistable.totalSupply({ from: admin });
        expect(result.toString()).to.equal(mintAmount.toString());
    });

    it('non-whitelisted cant recieve tokens from transfer', async function () {
        await this.ERC20Whitelistable.grantRole(MINTER_ROLE, authorized, { from: admin });
        await this.ERC20Whitelistable.grantRole(WHITELISTED_ROLE, authorized, { from: admin });
        await this.ERC20Whitelistable.mint(authorized, mintAmount, { from: authorized });

        await expectRevert(
            this.ERC20Whitelistable.transfer(other, transferAmount, { from: authorized }),
            'ERC20Whitelistable: must be whitelisted to recieve tokens',
        );
    });

    it('transfer tokens', async function () {
        await this.ERC20Whitelistable.grantRole(MINTER_ROLE, authorized, { from: admin });
        await this.ERC20Whitelistable.grantRole(WHITELISTED_ROLE, authorized, { from: admin });
        await this.ERC20Whitelistable.grantRole(WHITELISTED_ROLE, other, { from: admin });
        await this.ERC20Whitelistable.mint(authorized, mintAmount, { from: authorized });

        var receipt = await this.ERC20Whitelistable.transfer(other, transferAmount, { from: authorized })
        await expectEvent(receipt, "Transfer", { from: authorized, to: other, value: transferAmount });

        var result = await this.ERC20Whitelistable.balanceOf(other, { from: admin })
        expect(result.toString()).to.equal(transferAmount.toString());

        result = await this.ERC20Whitelistable.balanceOf(authorized, { from: admin })
        expect(result.toString()).to.equal(transferAmount.toString());

        result = await this.ERC20Whitelistable.totalSupply({ from: admin });
        expect(result.toString()).to.equal(mintAmount.toString());
    });
  });
});